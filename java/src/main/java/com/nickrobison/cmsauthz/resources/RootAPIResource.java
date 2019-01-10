package com.nickrobison.cmsauthz.resources;

import com.codahale.metrics.annotation.Timed;
import com.codahale.xsalsa20poly1305.SecretBox;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.type.TypeFactory;
import com.github.nitram509.jmacaroons.*;
import com.nickrobison.cmsauthz.api.JWKResponse;
import com.nickrobison.cmsauthz.helpers.Helpers;
import com.nickrobison.cmsauthz.helpers.VarInt;
import org.whispersystems.curve25519.Curve25519;
import org.whispersystems.curve25519.Curve25519KeyPair;

import javax.servlet.http.HttpServletRequest;
import javax.ws.rs.*;
import javax.ws.rs.client.Client;
import javax.ws.rs.core.*;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.security.SecureRandom;
import java.util.*;

@Path("/")
@Produces(MediaType.APPLICATION_JSON)
public class RootAPIResource {

    private static final Charset KEY_CHARSET = StandardCharsets.US_ASCII;
    private static final Charset MSG_CHARSET = StandardCharsets.UTF_8;
    private static final Base64.Decoder URL_DECODER = Base64.getUrlDecoder();
    private static final String TEST_NONCE = "this is a test nonce,...";
    private static String TEST_KEY = "this is a test key, it should be long enough.";

    private final Client client;
    private final String dischargeHost;


    public RootAPIResource(Client client) {
        System.out.println("Created`");
        this.client = client;
        final String host = System.getenv("HOST");
        if (host == null) {
            this.dischargeHost = "http://localhost:8080";
        } else {
            this.dischargeHost = host;
        }
    }

    @GET
    @Path("/{aco_id}/token")
    public Response getToken(@QueryParam("user_id") String userID, @PathParam("aco_id") String acoID, @QueryParam("vendor_id") Optional<String> vendorID) {

        // Get the JWKS
        final JWKResponse response = this.client
                .target(String.format("%s/api/acos/%s/.well-known/jwks.json", this.dischargeHost, acoID))
                .request(MediaType.APPLICATION_JSON_TYPE)
                .get(JWKResponse.class);

        // Decode the Key from Base64
        final byte[] decodedKey = URL_DECODER.decode(response.getKey());
        Helpers.printUnsignedBytes("Their pub key:", decodedKey);

        // We need to add a third party caveat to have the ACO endpoint give us a public key.

        // I think this is how it works.

//        Create a new keypair
        final Curve25519 cipher = Curve25519.getInstance(Curve25519.BEST);
        final Curve25519KeyPair keyPair = cipher.generateKeyPair();

        // Nonce
        final SecureRandom secureRandom = new SecureRandom();
        byte[] nonce = new byte[24];
        secureRandom.nextBytes(nonce);

        Helpers.printUnsignedBytes("Nonce", TEST_NONCE.getBytes());
        Helpers.printUnsignedBytes("Pub key", keyPair.getPublicKey());

        // Encrypt things
        final String caveat = String.format("user_id= %s", userID);
        final byte[] encrypted = encodeIdentifier(keyPair, decodedKey, TEST_KEY, TEST_NONCE.getBytes(), caveat);

        final Macaroon m2 = new MacaroonsBuilder("http://localhost:3002/", TEST_KEY, "first-party-id", MacaroonVersion.VERSION_2)
//                .add_first_party_caveat(String.format("aco_id= %s", acoID))
                .add_third_party_caveat(String.format("%s/api/acos/%s/verify", this.dischargeHost, acoID), TEST_KEY, encrypted)
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

        MacaroonsVerifier verifier = new MacaroonsVerifier(macaroons.get().get(0));

        for (Macaroon discharged : macaroons.get().subList(1, macaroons.get().size())) {
            verifier = verifier.satisfy3rdParty(discharged);
        }

//        Add all the discharges
//        final MacaroonsVerifier verifier = macaroons.get().subList(1, macaroons.get().size())
//                .stream()
//                .reduce(initialVerifier, MacaroonsVerifier::satisfy3rdParty, (acc, mac) -> acc);

        boolean valid;

        try {
            verifier
                    .satisfyExact(String.format("aco_id= test_aco"))
                    .assertIsValid(TEST_KEY);
            valid = true;
        } catch (GeneralSecurityRuntimeException | MacaroonValidationException e) {
            System.out.println(e.getMessage());
            valid = false;
        }

        if (valid) {
            return Response.ok(String.format("Successfully accessed data for ACO %s", aco_id)).build();
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
                .map(RootAPIResource::parseMacaroons);
//                .map(macs -> macs.stream().map(MacaroonsBuilder::deserialize).collect(Collectors.toList()));
    }

    private static List<Macaroon> parseMacaroons(String macaroon) {

        //        Figure out what we're looking at.
        switch (macaroon.charAt(0)) {
//            JSON v2
            case '{': {
                throw new IllegalArgumentException("Cannot decode V2 JSON");
            }
            // V2 Binary
            case 0x02: {
                throw new IllegalArgumentException("Cannot decode V2 Binary");
            }
//            If we're an array, pull out everything as a string
            case '[': {
                if (macaroon.charAt(1) == '{') {

                    // We're an array of V2 Macaroons, so we just need to split the array
                    // Decode an array of V2 Macaroons
                    final ObjectMapper objectMapper = new ObjectMapper();
                    final TypeFactory typeFactory = objectMapper.getTypeFactory();
                    try {
                        final List<Object> read = objectMapper.readValue(macaroon, typeFactory.constructCollectionType(List.class, Object.class));
                        return Collections.emptyList();
                    } catch (IOException e) {
                        throw new IllegalStateException(e);
                    }
                }
                break;
            }
            // If 0 is the first thing, we're just a single Macaroon.
            // That means we can pass the base64 encoded string directly to the Macaroons library
            case '0': {
                return Collections.singletonList(MacaroonsBuilder.deserialize(macaroon));
            }
        }

        // Decode it, and do it all over again
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
//            If we're an array, pull out everything as a string
            case '[': {
                if (decoded[1] == '{') {

                    // We're an array of V2 Macaroons, so we just need to split the array
                    // Decode an array of V2 Macaroons
                    final ObjectMapper objectMapper = new ObjectMapper();
                    List<Macaroon> macaroons = new ArrayList<>();

                    try {
                        final JsonNode tree = objectMapper.readTree(decoded);
                        for (final JsonNode node : tree) {
                            final String nodeText = node.toString();
                            macaroons.add(MacaroonsBuilder.deserialize(nodeText));
                        }
                        return macaroons;
                    } catch (IOException e) {
                        throw new IllegalStateException(e);
                    }
                }
                break;
            }
            // If 0 is the first thing, we're just a single Macaroon.
            // That means we can pass the base64 encoded string directly to the Macaroons library
            case '0': {
                return Collections.singletonList(MacaroonsBuilder.deserialize(new String(decoded)));
            }
        }

        // Now we know that we think we're an array of V1 Macaroons, make them a json array and split them apart.
        final ObjectMapper objectMapper = new ObjectMapper();
        final TypeFactory typeFactory = objectMapper.getTypeFactory();
        try {
            return objectMapper.readValue(macaroon, typeFactory.constructCollectionType(List.class, String.class));
        } catch (IOException e) {
            throw new IllegalStateException(e);
        }
    }

    /**
     * Encode a caveat identifier using the V2 serialization format.
     * The message (along with the root key) is encrypted as a NaCL secret box, using the {@link SecretBox#seal(byte[], byte[])} method.
     * The encrypted method is then appended to the end of the caveat ID, which contains the public key used to encrypt the data (the public key of this service),
     * 4 bytes which identify the public key (of the third party service) used, and the nonce.
     * <p>
     * The returned string is Base64 encoded using the {@link Base64#getEncoder()} method. Libmacaroons expects caveats to have non-URL safe encoding.
     * So we acquiesce.
     *
     * @param keyPair             - {@link Curve25519KeyPair} of this service, which is used to encrypt the message as a {@link SecretBox}
     * @param thirdPartyPublicKey - {@link byte[]} 32-byte Curve-25519 public key of the third party service
     * @param rootKey             - {@link String} root key used with the message.
     * @param nonce               - {@link byte[]} 24-byte random nonce used in encryption
     * @param message             - {@link String} secret message to encrypt and send to third-party
     * @return - {@link String} Base64 encoded Caveat ID
     */
    private static byte[] encodeIdentifier(Curve25519KeyPair keyPair, byte[] thirdPartyPublicKey, String rootKey, byte[] nonce, String message) {

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

        // We need to use the ISO 8859 charset in order to keep the byte values the same.
        // See: https://stackoverflow.com/a/17575008
        return fullMessage.array();

//        return Base64.getEncoder().encodeToString(fullMessage.array());
    }
}
