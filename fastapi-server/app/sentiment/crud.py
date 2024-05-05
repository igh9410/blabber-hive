# app/sentiment/crud.py
"""
The crud module contains the business logic for the sentiment analysis service.
It defines a single function, analyze_sentiment,
which takes a text string as input
and returns the sentiment of
the text as a string.
The function uses the SentimentModel
to predict the sentiment of the given text.
"""
from app.models.sentiment_model import SentimentModel

sentiment_model = SentimentModel()


def analyze_sentiment(text: str) -> str:
    """
    Analyze the sentiment of the given text using the sentiment analysis model.

    Args:
        text (str): The text to analyze.

    Returns:
        str: The predicted sentiment of the text.
    """
    sentiment = sentiment_model.predict_sentiment(text)
    return sentiment
