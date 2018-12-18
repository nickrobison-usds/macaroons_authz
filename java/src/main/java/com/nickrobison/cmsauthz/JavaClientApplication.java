package com.nickrobison.cmsauthz;

import com.nickrobison.cmsauthz.resources.RootAPIResource;
import io.dropwizard.Application;
import io.dropwizard.client.JerseyClientBuilder;
import io.dropwizard.setup.Bootstrap;
import io.dropwizard.setup.Environment;

import javax.ws.rs.client.Client;

public class JavaClientApplication extends Application<JavaClientConfiguration> {

    public static void main(final String[] args) throws Exception {
        new JavaClientApplication().run(args);
    }

    @Override
    public String getName() {
        return "JavaClient";
    }

    @Override
    public void initialize(final Bootstrap<JavaClientConfiguration> bootstrap) {
        // TODO: application initialization
    }

    @Override
    public void run(final JavaClientConfiguration configuration,
                    final Environment environment) {

        final Client client = new JerseyClientBuilder(environment)
                .using(configuration.getHttpClient())
                .build(getName());
        environment.jersey().register(new RootAPIResource(client));
    }
}
