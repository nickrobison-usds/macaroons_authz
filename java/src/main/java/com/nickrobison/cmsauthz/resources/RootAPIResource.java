package com.nickrobison.cmsauthz.resources;

import com.codahale.metrics.annotation.Timed;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.type.TypeFactory;
import com.github.nitram509.jmacaroons.Macaroon;
import com.github.nitram509.jmacaroons.MacaroonsBuilder;
import com.github.nitram509.jmacaroons.MacaroonsVerifier;
import com.neilalexander.jnacl.NaCl;
import com.nickrobison.cmsauthz.api.JWKResponse;
import org.apache.commons.codec.binary.Base64;

import javax.servlet.http.HttpServletRequest;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.client.Client;
import javax.ws.rs.core.*;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.SecureRandom;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

@Path("/")
@Produces(MediaType.APPLICATION_JSON)
public class RootAPIResource {

    private final Map<String, String> keyMap;
    private final Client client;
    private final String TEST_KEY;

    public RootAPIResource(Client client) {
        this.client = client;
        this.keyMap = new ConcurrentHashMap<>();
        final SecureRandom secureRandom = new SecureRandom();
        byte[] bytes = new byte[32];
        secureRandom.nextBytes(bytes);
        this.TEST_KEY = new String(bytes, StandardCharsets.UTF_8);
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
        final byte[] decodedKey = Base64.decodeBase64(response.getKey());

        // We need to add a third party caveat to have the ACO endpoint give us a public key. 
//        final Macaroon macaroon = new MacaroonsBuilder("http://localhost:3002/", TEST_KEY, "first-party-id")
//                .add_first_party_caveat("aco_id = 1")
//                .getMacaroon();

        // I think this is how it works.

//        Create a new keypair
        final KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        final KeyPair keyPair = generator.generateKeyPair();


        // Encrypt things

        // Create a new secret box using our private key and their public key
        final NaCl naCl = new NaCl(keyPair.getPrivate().getEncoded(), decodedKey);

        // Create a 24 byte long random nonce
        final SecureRandom secureRandom = new SecureRandom();
        final byte[] nonce = new byte[24];
        secureRandom.nextBytes(nonce);

        // Encrypt it
        final String msg = "This is a test message, does it work?";
        final byte[] encrypted = naCl.encrypt(msg.getBytes(), nonce);

        final Macaroon m2 = new MacaroonsBuilder("http://localhost:3002/", TEST_KEY, "first-party-id")
                .add_third_party_caveat("http://localhost:8080/api/users/verify", TEST_KEY, new String(encrypted, StandardCharsets.UTF_8))
                .getMacaroon();

        return Response.ok().entity(m2.serialize()).build();
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
                .satisfyExcact("aco_id = 1")
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
        final Base64 base64 = new Base64();
        final byte[] decoded = base64.decode(macaroon);
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
}
