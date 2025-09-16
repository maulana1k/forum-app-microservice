from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

# Initialize router
predict_router = APIRouter()

# Define request body
class ToxicityRequest(BaseModel):
    postId: int
    content: str

# Initialize Hugging Face pipeline (runs once at startup)
# classifier = pipeline("text-classification", model="unitary/toxic-bert", return_all_scores=True)
# pipe = pipeline("text-classification", model="ayameRushia/bert-base-indonesian-1.5G-sentiment-analysis-smsa")
# pipe = pipeline("text-classification", model="ayameRushia/roberta-base-indonesian-sentiment-analysis-smsa", top_k=None)

@predict_router.post("/predict/toxicity", tags=['Toxicity'])
def predict_toxicity(req: ToxicityRequest):
    """
    Predict toxicity level of a post's content
    """
    if not req.content.strip():
        raise HTTPException(status_code=400, detail="Content cannot be empty")

    # Run model
    results = {"sample":0.9124} # returns list of dicts for each label
    return results
    # Convert to dict: label -> score
    # toxicity_scores = {r['label']: r['score'] for r in results}

    # # Determine overall toxicity (optional: pick label with highest score)
    # max_label = max(results, key=lambda x: x['score'])['label']

    # return {
    #     "postId": req.postId,
    #     "content": req.content,
    #     "toxicity_scores": toxicity_scores,
    #     "predicted_label": max_label
    # }
