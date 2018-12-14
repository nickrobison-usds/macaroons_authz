package com.nickrobison.cmsauthz.resources;

import com.codahale.metrics.annotation.Timed;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

@Path("/")
@Produces(MediaType.APPLICATION_JSON)
public class RootAPIResource {

    public RootAPIResource() {
//        Not used
    }

    @GET
    @Timed
    public Response get() {
        return Response.ok("Hello there!").build();
    }
}
