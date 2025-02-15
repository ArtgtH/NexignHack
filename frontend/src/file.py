import requests
import streamlit as st
import pandas as pd
import plotly.express as px


RESULT = {0: "NEUTRAL", 1: "POSITIVE", 2: "NEGATIVE"}


def send_file_to_backend(file, url):
    files = {"file": file}
    res = requests.post(url, files=files)
    return res.json()


def download_file(url):
    st.text("Загрузка Excel-файла")

    uploaded_file = st.file_uploader("Выберите Excel-файл", type=["xlsx", "xls"])

    if st.button("Отправить файл"):
        if uploaded_file is not None:
            res = send_file_to_backend(uploaded_file, url)

            df = pd.DataFrame(res["messages"])
            df["result"] = df["result"].map(RESULT)

            st.write(df)

            st.subheader("Распределение результатов")

            result_counts = df["result"].value_counts().reset_index()
            result_counts.columns = ["Category", "Count"]

            colors = {
                "POSITIVE": "#4CAF50",
                "NEUTRAL": "#FFC107",
                "NEGATIVE": "#F44336",
            }

            fig = px.bar(
                result_counts,
                x="Category",
                y="Count",
                color="Category",
                color_discrete_map=colors,
                text="Count",
                height=400,
            )

            fig.update_layout(
                title_text="Распределение эмоциональной окраски",
                title_x=0.5,
                xaxis_title="Категория",
                yaxis_title="Количество",
                showlegend=False,
                margin=dict(l=20, r=20, t=60, b=20),
            )

            fig.update_traces(
                texttemplate="%{text}",
                textposition="outside",
                marker_line_color="rgb(8,48,107)",
                marker_line_width=1.5,
            )

            st.plotly_chart(fig, use_container_width=True)

        else:
            st.warning("Пожалуйста, загрузите файл")
