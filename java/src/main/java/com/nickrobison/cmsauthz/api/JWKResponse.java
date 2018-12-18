package com.nickrobison.cmsauthz.api;

import com.fasterxml.jackson.annotation.JsonProperty;

public class JWKResponse {

    private String algorithm;
    private String id;
    private String use;
    private String key;

    @JsonProperty("alg")
    public String getAlgorithm() {
        return algorithm;
    }

    @JsonProperty("alg")
    public void setAlgorithm(String algorithm) {
        this.algorithm = algorithm;
    }

    @JsonProperty("kid")
    public String getId() {
        return id;
    }

    @JsonProperty("kid")
    public void setId(String id) {
        this.id = id;
    }

    public String getUse() {
        return use;
    }

    public void setUse(String use) {
        this.use = use;
    }

    @JsonProperty("k")
    public String getKey() {
        return key;
    }

    @JsonProperty("k")
    public void setKey(String key) {
        this.key = key;
    }
}
