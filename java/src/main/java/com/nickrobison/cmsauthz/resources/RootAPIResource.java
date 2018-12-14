package com.nickrobison.cmsauthz.resources;

import com.codahale.metrics.annotation.Timed;
import com.github.nitram509.jmacaroons.Macaroon;
import com.github.nitram509.jmacaroons.MacaroonsBuilder;
import com.github.nitram509.jmacaroons.MacaroonsVerifier;

import javax.servlet.http.HttpServletRequest;
import javax.ws.rs.*;
import javax.ws.rs.core.*;
import java.util.Map;
import java.util.Optional;

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

        final Optional<Macaroon> macaroon = this.getMacaroonFromHeader(headers);

        if (!macaroon.isPresent()) {
            return Response.status(Response.Status.UNAUTHORIZED).entity("Must have Macaroon").build();
        }

        final MacaroonsVerifier verifier = new MacaroonsVerifier(macaroon.get())
                .satisfyExcact("aco_id = 1");
        final boolean valid = verifier.isValid(TEST_KEY);

        if (valid) {
            return Response.ok("Hello there!").build();
        }

        return Response.status(Response.Status.UNAUTHORIZED).entity("Incorrect macaroon").build();
    }


    /**
     * Get Macaroon from header.
     * If it's missing, that should fail the auth.
     * @param headers - {@link HttpHeaders} headers from request
     * @return - {@link Optional} {@link Macaroon} from header
     */
    private Optional<Macaroon> getMacaroonFromHeader(HttpHeaders headers) {
        final Map<String, Cookie> cookies = headers.getCookies();
        Optional<Cookie> firstMacaroon = cookies
                .entrySet()
                .stream()
                .filter((entry) -> entry.getKey().startsWith("macaroon-"))
                .map(Map.Entry::getValue)
                .findFirst();

        return firstMacaroon.map(cookie -> MacaroonsBuilder.deserialize(cookie.getValue()));
    }
}
