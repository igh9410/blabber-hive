"""
transformers is a library that provides a wide variety of pre-trained models
for Natural Language Processing (NLP) tasks.
"""

from transformers import AutoModelForSequenceClassification, AutoTokenizer


class SentimentModel:
    """
    SentimentModel class is a wrapper around
    a pre-trained model for sentiment analysis.
    """

    def __init__(
        self,
        model_name="distilbert-base-uncased-finetuned-sst-2-english",
    ):
        self.model_name = model_name
        self.model = AutoModelForSequenceClassification.from_pretrained(
            model_name
        )
        self.tokenizer = AutoTokenizer.from_pretrained(model_name)

    def predict_sentiment(self, text):
        """
        predict_sentiment method takes a text input and returns
        the predicted sentiment.
        """
        inputs = self.tokenizer(
            text, return_tensors="pt", padding=True, truncation=True
        )
        outputs = self.model(**inputs)
        predicted_class = outputs.logits.argmax().item()

        if predicted_class == 0:
            return "negative"
        elif predicted_class == 1:
            return "neutral"
        else:
            return "positive"
