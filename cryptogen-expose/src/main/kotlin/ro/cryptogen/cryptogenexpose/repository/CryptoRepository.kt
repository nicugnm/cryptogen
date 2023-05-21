package ro.cryptogen.cryptogenexpose.repository

import org.springframework.data.mongodb.repository.MongoRepository
import ro.cryptogen.cryptogenexpose.model.CryptoMessage

interface CryptoRepository : MongoRepository<CryptoMessage, String> {
}