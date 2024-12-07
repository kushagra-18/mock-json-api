package com.mock_json.mock_api.repositories;


import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

import com.mock_json.mock_api.models.RandomWords;

@Repository
public interface RandomWordsRepository extends MongoRepository<RandomWords, String> {


}
