import requests
import streamlit as st


def send_file_to_backend(file, url):
	files = {'file': file}
	res = requests.post(url, files=files)
	return res


def download_file(url):
	st.text("Загрузка Excel-файла")

	uploaded_file = st.file_uploader("Выберите Excel-файл", type=["xlsx", "xls"])

	if st.button("Отправить файл"):
		if uploaded_file is not None:
			res = send_file_to_backend(uploaded_file, url)

			if res.status_code == 201:
				st.success("Файл успешно отправлен!")
				st.text(res.text)
			else:
				st.error(f"Ошибка при отправке файла: {res.status_code}")
		else:
			st.warning("Пожалуйста, загрузите файл и укажите URL бэкенда.")
