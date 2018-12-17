package com.nickrobison.cmsauthz.resources;

import com.codahale.metrics.annotation.Timed;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.type.TypeFactory;
import com.github.nitram509.jmacaroons.Macaroon;
import com.github.nitram509.jmacaroons.MacaroonsBuilder;
import com.github.nitram509.jmacaroons.MacaroonsVerifier;
import org.apache.commons.codec.binary.Base64;

import javax.servlet.http.HttpServletRequest;
import javax.ws.rs.*;
import javax.ws.rs.core.*;
import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

@Path("/")
@Produces(MediaType.APPLICATION_JSON)
public class RootAPIResource {

    private static final String TEST_KEY = "test key";

    public RootAPIResource() {
//        Not used
    }

    @GET
    @Path("/token")
    public Response getToken() {
        final Macaroon macaroon = new MacaroonsBuilder("http://localhost:3002/", TEST_KEY, "test-token-1")
                .add_first_party_caveat("aco_id = 1")
                .getMacaroon();

        return Response.ok().entity(macaroon.serialize()).build();
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
