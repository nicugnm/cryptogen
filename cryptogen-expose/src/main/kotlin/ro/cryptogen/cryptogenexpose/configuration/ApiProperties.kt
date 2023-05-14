package ro.cryptogen.cryptogenexpose.configuration

import org.springframework.boot.context.properties.ConfigurationProperties

@ConfigurationProperties(prefix = "api")
class ApiProperties {
    lateinit var url: String
}