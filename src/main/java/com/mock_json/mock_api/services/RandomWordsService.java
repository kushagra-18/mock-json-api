package com.mock_json.mock_api.services;


import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.data.mongodb.core.aggregation.Aggregation;
import org.springframework.data.mongodb.core.aggregation.AggregationResults;
import org.springframework.stereotype.Service;

import com.mock_json.mock_api.models.RandomWords;

import java.util.List;

@Service
public class RandomWordsService {

    private final MongoTemplate mongoTemplate;

    public RandomWordsService(MongoTemplate mongoTemplate) {
        this.mongoTemplate = mongoTemplate;
    }

    public String getRandomSlug() {
        List<RandomWords> randomWords = getRandomWords(2);
        return randomWords.get(0).getWord() + "-" + randomWords.get(1).getWord();
    }

    private List<RandomWords> getRandomWords(int count) {
        Aggregation aggregation = Aggregation.newAggregation(
                Aggregation.sample(count));

        AggregationResults<RandomWords> results = mongoTemplate.aggregate(aggregation, "random_words",
                RandomWords.class);
        return results.getMappedResults();
    }
}
