import os

import streamlit as st

from file import download_file
from text import download_text


def main():
	back = os.getenv("BACKEND_URL")

	st.title("AI Learning Lab")

	st.sidebar.title("Меню")
	page = st.sidebar.selectbox(
		"Выберите страницу",
		("Файл", "Строка"),
		index=0
	)

	if page == "Файл":
		download_file(back + "full/")

	else:
		download_text(back + "short/")


if __name__ == "__main__":
	main()
