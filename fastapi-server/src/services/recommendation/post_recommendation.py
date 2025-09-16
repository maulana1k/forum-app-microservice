import grpc
from concurrent import futures
import logging
from typing import List

from src.grpc import recommender_pb2
from src.grpc import recommender_pb2_grpc

# from src.ml.models.predictor import MLModelManager  # Your ML model

class RecommenderService(recommender_pb2_grpc.RecommenderServiceServicer):
    async def GetRecommendedPosts(self, request, context):
        try:
            # Your recommendation logic here
            # This is where you'd call your ML model
            recommended_posts = await self._get_recommendations(
                request.user_id, 
                request.topic, 
                request.limit
            )
            
            return recommender_pb2.RecommendationResponse(posts=recommended_posts)
            
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error getting recommendations: {str(e)}")
            return recommender_pb2.RecommendationResponse()

    async def _get_recommendations(self, user_id: str, topic: str, limit: int) -> List[recommender_pb2.PostItem]:
        # Example implementation - replace with your actual logic
        # This would typically call your ML model
        
        # For now, return mock data
        posts = []
        for i in range(min(limit, 5)):  # Mock 5 posts max
            posts.append(recommender_pb2.PostItem(
                post_id=f"post_{user_id}_{i}",
                content=f"Recommended content about {topic} for user {user_id}",
                author_id=f"author_{i}",
            ))
        
        return posts