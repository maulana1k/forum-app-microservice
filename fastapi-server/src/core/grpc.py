import asyncio
import grpc
from concurrent import futures
import logging
from typing import Optional

from src.grpc import recommender_pb2_grpc
from src.services.recommendation.post_recommendation import RecommenderService
from src.core.interceptor import LoggingInterceptor, DetailedLoggingInterceptor


class GRPCServer:
    def __init__(self, host: str = "0.0.0.0", port: int = 50051):
        self.host = host
        self.port = port
        self.server: Optional[grpc.aio.Server] = None
        self._server_task: Optional[asyncio.Task] = None

    async def start(self):
        """Start the gRPC server"""
        if self.server is not None:
            return
        interceptors = [LoggingInterceptor(), DetailedLoggingInterceptor()]

        self.server = grpc.aio.server(
            futures.ThreadPoolExecutor(max_workers=10), interceptors=interceptors
        )

        # Add services
        recommender_pb2_grpc.add_RecommenderServiceServicer_to_server(
            RecommenderService(), self.server
        )

        # Listen on port
        listen_addr = f"{self.host}:{self.port}"
        self.server.add_insecure_port(listen_addr)

        # Start server
        await self.server.start()
        logging.info(f"gRPC server started on {listen_addr}")

        # Create a task to keep server running
        self._server_task = asyncio.create_task(self._serve())

    async def _serve(self):
        """Keep the server running"""
        if self.server:
            await self.server.wait_for_termination()

    async def stop(self, grace: int = 1):
        """Stop the gRPC server gracefully."""
        if not self.server:
            return

        # Cancel the serve task if running
        if self._server_task:
            self._server_task.cancel()
            try:
                await self._server_task
            except asyncio.CancelledError:
                pass

        # Stop the gRPC server with a grace period
        await self.server.stop(grace)
        self.server = None
        logging.info("gRPC server stopped")


# Global gRPC server instance
grpc_server = GRPCServer()
