from flask import Flask, render_template

import requests
api_url = "http://localhost:8080"


def getAllArticles(url):
    """Get a list with all articles in json format."""
    response = requests.get(f'{url}/articles')
    articles = response.json()
    return articles


def DeleteArticleById(url, id):
    """Delete an article by Id."""
    response = requests.get(f'{url}/article/delete?id={id}')
    articles = response.json()
    return articles

app = Flask(__name__)

@app.route('/')
def hello():
    return 'Hello, World!'


@app.route('/articles')
def index():
    articles = getAllArticles(api_url)
    return render_template('index.html', articles=articles)
