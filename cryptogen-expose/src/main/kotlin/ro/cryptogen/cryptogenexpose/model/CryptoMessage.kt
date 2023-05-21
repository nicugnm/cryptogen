package ro.cryptogen.cryptogenexpose.model

import org.springframework.data.mongodb.core.mapping.Document

@Document(collection = "crypto_messages")
data class CryptoMessage(
    val content: String
)