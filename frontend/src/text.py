import pandas as pd
import requests
import streamlit as st

from file import RESULT


def send_text_to_backend(text, url):
    data = {"text": text}
    res = requests.post(url, json=data)
    return res.json()


def download_text(url):
    st.text("Отправка построчно")

    phrase = st.text_input("Введите фразу: ")

    if st.button("Отправить текст") and phrase:
        res = send_text_to_backend(phrase, url)

        st.write(RESULT[res["result"]])
