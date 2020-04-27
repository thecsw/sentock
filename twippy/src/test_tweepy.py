import pytest
import tweepy
import random
import os
from twippy import got_tweet
from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer

def test_fake_tweets():
    companies = [
        "McDonalds",
        "Chipotle",
        "Microsoft",
        "Fedex",
        "Disney",
        "Google",
        "Tesla",
        "Facebook",
        "Amazon",
    ]
    test_tweets = [
        "McDonald's is sooo tasty!",
        "I freaking hate Chipotle!",
        "I do not dislike Microsoft",
        "Fedex is te wey bruddah *clicking intensifies*",
        "I love Disney!!!!!",
        "Google hates me, and I hate Google",
        "Elon Musk is doing a great job at Tesla",
        "Mark Zuckerberg's Facebook is horrible",
        "I love Amazon's customer service!",
    ]

    r = random.Random()

    tested_tweets = []
    for i in range(0, 9):
        tested = got_tweet(
            r.randrange(1, 4294967295), test_tweets[i], r.randrange(1, 4294967295), True
        )
        assert tested == companies[i] or tested == "No sentiment value"


def test_sentiment_values():
    analyzer = SentimentIntensityAnalyzer()
    test_tweets = [
        "McDonald's is sooo tasty!",
        "I freaking hate Chipotle!",
        "I do not dislike Microsoft",
        "I love Disney!!!!!",
        "Google hates me, and I hate Google",
        "Elon Musk is doing a great job at Tesla",
        "Mark Zuckerberg's Facebook is horrible",
        "Amazon's customer service is stellar!",
    ]
    positive_test_tweets = [
        "McDonald's is sooo tasty!",
        "I do not dislike Microsoft",
        "I love Disney!!!!!",
        "Elon Musk is doing a great job at Tesla",
        "Amazon's customer service is stellar!",
    ]
    negative_test_tweets = [
        "I freaking hate Chipotle!",
        "Google hates me, and I hate Google",
        "Mark Zuckerberg's Facebook is horrible",
    ]
    for tweet in positive_test_tweets:
        assert (
            float(analyzer.polarity_scores(tweet)["compound"]) > 0
            or float(analyzer.polarity_scores(tweet)["compound"]) == 0
        )
    for tweet in negative_test_tweets:
        assert (
            float(analyzer.polarity_scores(tweet)["compound"]) < 0
            or float(analyzer.polarity_scores(tweet)["compound"]) == 0
        )
