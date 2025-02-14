from pydantic import Field
from pydantic_settings import BaseSettings


class AppConfig(BaseSettings):
    backend_url: str = Field(env='BACKEND_URL')

    class Config:
        env_file = ".env"
        env_file_encoding = 'utf-8'


