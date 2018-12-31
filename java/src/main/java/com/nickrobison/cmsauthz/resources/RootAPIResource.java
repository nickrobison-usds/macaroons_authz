package com.nickrobison.cmsauthz.resources;

import com.codahale.metrics.annotation.Timed;
import com.codahale.xsalsa20poly1305.SecretBox;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.type.TypeFactory;
import com.github.nitram509.jmacaroons.Macaroon;
import com.github.nitram509.jmacaroons.MacaroonVersion;
import com.github.nitram509.jmacaroons.MacaroonsBuilder;
import com.github.nitram509.jmacaroons.MacaroonsVerifier;
import com.nickrobison.cmsauthz.api.JWKResponse;
import com.nickrobison.cmsauthz.helpers.VarInt;
import org.whispersystems.curve25519.Curve25519;
import org.whispersystems.curve25519.Curve25519KeyPair;

import javax.servlet.http.HttpServletRequest;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.client.Client;
import javax.ws.rs.core.*;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.security.SecureRandom;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

@Path("/")
@Produces(MediaType.APPLICATION_JSON)
public class RootAPIResource {

    private static final Charset KEY_CHARSET = StandardCharsets.US_ASCII;
    private static final Charset MSG_CHARSET = StandardCharsets.UTF_8;
    public static final Base64.Encoder URL_ENCODER = Base64.getUrlEncoder();
    public static final Base64.Decoder URL_DECODER = Base64.getUrlDecoder();

    private final Map<String, String> keyMap;
    private final Client client;
    private final String TEST_KEY;

    public RootAPIResource(Client client) {
        this.client = client;
        this.keyMap = new ConcurrentHashMap<>();
        final SecureRandom secureRandom = new SecureRandom();
        byte[] bytes = new byte[32];
        secureRandom.nextBytes(bytes);
        this.TEST_KEY = new String(bytes, MSG_CHARSET);
    }

    @GET
    @Path("/token")
    public Response getToken() throws Exception {

        // Get the JWKS
        final JWKResponse response = this.client
                .target("http://localhost:8080/api/users/.well-known/jwks.json")
                .request(MediaType.APPLICATION_JSON_TYPE)
                .get(JWKResponse.class);

        // Decode the Key from Base64
        final byte[] decodedKey = URL_DECODER.decode(response.getKey());
        printUnsignedBytes("Their pub key:", decodedKey);

        // We need to add a third party caveat to have the ACO endpoint give us a public key.

        // I think this is how it works.

//        Create a new keypair
        final Curve25519 cipher = Curve25519.getInstance(Curve25519.BEST);
        final Curve25519KeyPair keyPair = cipher.generateKeyPair();

        // Nonce
        final SecureRandom secureRandom = new SecureRandom();
        byte[] nonce = new byte[24];
        secureRandom.nextBytes(nonce);

        printUnsignedBytes("Nonce", nonce);
        printUnsignedBytes("Pub key", keyPair.getPublicKey());

        // Encrypt things
        final String tkey = "this is a test key, it should be long enough.";
        final String encrypted = encodeIdentifier(keyPair, decodedKey, tkey, nonce, "this is a test message");

        printUnsignedBytes("Full ID", Base64.getDecoder().decode(encrypted));

        final Macaroon m2 = new MacaroonsBuilder("http://localhost:3002/", TEST_KEY, "first-party-id", MacaroonVersion.VERSION_1)
                .add_third_party_caveat("http://localhost:8080/api/users/verify", tkey, encrypted)
                .getMacaroon();

        return Response.ok().entity(m2.serialize(MacaroonVersion.SerializationVersion.V2_JSON)).build();
    }

    @GET
    @Path("/{aco_id}")
    @Timed
    public Response get(@Context HttpServletRequest request, @PathParam("aco_id") String aco_id, @Context HttpHeaders headers) {

        final Optional<List<Macaroon>> macaroons = this.getMacaroonFromHeader(headers);

        if (!macaroons.isPresent()) {
            return Response.status(Response.Status.UNAUTHORIZED).entity("Must have Macaroon").build();
        }

        final MacaroonsVerifier initialVerifier = new MacaroonsVerifier(macaroons.get().get(0));

//        Add all the discharges
        final MacaroonsVerifier verifier = macaroons.get().subList(1, macaroons.get().size())
                .stream()
                .reduce(initialVerifier, MacaroonsVerifier::satisfy3rdParty, (acc, mac) -> acc);

        final boolean valid = verifier
                .satisfyExact("aco_id = 1")
                .isValid(TEST_KEY);

        if (valid) {
            return Response.ok("Hello there!").build();
        }

        return Response.status(Response.Status.UNAUTHORIZED).entity("Incorrect macaroon").build();
    }


    /**
     * Get Macaroon from header.
     * If it's missing, that should fail the auth.
     *
     * @param headers - {@link HttpHeaders} headers from request
     * @return - {@link Optional} {@link Macaroon} from header
     */
    private Optional<List<Macaroon>> getMacaroonFromHeader(HttpHeaders headers) {
        final Map<String, Cookie> cookies = headers.getCookies();
        Optional<Cookie> firstMacaroon = cookies
                .entrySet()
                .stream()
                .filter((entry) -> entry.getKey().startsWith("macaroon-"))
                .map(Map.Entry::getValue)
                .findFirst();

        // Base64 decode the input
        return firstMacaroon
                .map(Cookie::getValue)
                .map(RootAPIResource::parseMacaroons)
                .map(macs -> macs.stream().map(MacaroonsBuilder::deserialize).collect(Collectors.toList()));
    }

    private static List<String> parseMacaroons(String macaroon) {
        final byte[] decoded = URL_DECODER.decode(macaroon);
//        Figure out what we're looking at.
        switch (decoded[0]) {
//            JSON v2
            case '{': {
                throw new IllegalArgumentException("Cannot decode V2 JSON");
            }
            // V2 Binary
            case 0x02: {
                throw new IllegalArgumentException("Cannot decode V2 Binary");
            }
//            If we're an array, maybe do things.
            case '[': {
//                Assume we're an array of JSON, and bail
                if (decoded[1] == '{') {
                    throw new IllegalArgumentException("Cannot decode array of V2 JSON");
                }
            }
            // If 0 is the first thing, we're just a single Macaroon.
            // That means we can pass the base64 encoded string directly to the Macaroons library
            case '0': {
                return Collections.singletonList(macaroon);
            }
        }

        // Now we know that we're an array of V1 Macaroons, so make them a json array and split them apart.
        final ObjectMapper objectMapper = new ObjectMapper();
        final TypeFactory typeFactory = objectMapper.getTypeFactory();
        try {
            return objectMapper.readValue(decoded, typeFactory.constructCollectionType(List.class, String.class));
        } catch (IOException e) {
            throw new IllegalStateException(e);
        }
    }

    private static String encodeIdentifier(Curve25519KeyPair keyPair, byte[] thirdPartyPublicKey, String rootKey, byte[] nonce, String message) throws Exception {

        final byte[] keyBytes = rootKey.getBytes(KEY_CHARSET);
        final byte[] msgBytes = message.getBytes(MSG_CHARSET);

        // Create Varint of root key length
        final byte[] tBytes = VarInt.writeUnsignedVarInt(keyBytes.length);

        // Allocate a byte buffer that is the size of the version (1 byte), the varint length of rootKey, rootKey and message
        final ByteBuffer msgBuffer = ByteBuffer.allocate(1
                + tBytes.length
                + keyBytes.length
                + msgBytes.length);

        // Add everything
        msgBuffer.put((byte) 2);
        msgBuffer.put(tBytes);
        msgBuffer.put(keyBytes);
        msgBuffer.put(msgBytes);
        // Reset it
        msgBuffer.flip();

        final SecretBox sbox = new SecretBox(thirdPartyPublicKey, keyPair.getPrivateKey());

        final byte[] sealed = sbox.seal(nonce, msgBuffer.array());

//        Now, add the header
        final ByteBuffer fullMessage = ByteBuffer.allocate(1 + 4 + 32 + 24 + sealed.length);
        fullMessage.put((byte) 2);
        fullMessage.put(Arrays.copyOfRange(thirdPartyPublicKey, 0, 4));
        fullMessage.put(keyPair.getPublicKey());
        fullMessage.put(nonce);
        fullMessage.put(sealed);
        fullMessage.flip();

        return Base64.getEncoder().encodeToString(fullMessage.array());
    }

    private static void printUnsignedBytes(String name, byte[] signedBytes) {

        int[] unsigned = new int[signedBytes.length];

        for (int i = 0; i < signedBytes.length; i++) {
            unsigned[i] = (signedBytes[i] & 0xFF);
        }

        System.out.printf("%s as unsigned bytes:\n", name);
        System.out.println(Arrays.toString(unsigned));
    }
}
