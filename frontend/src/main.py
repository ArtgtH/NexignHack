import streamlit as st

from components import download_file
from config import AppConfig


def main():
    config = AppConfig()
    st.title('AI Learning Lab')
    download_file(config.backend_url)


if __name__ == '__main__':
    main()
