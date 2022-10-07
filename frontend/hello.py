from flask import Flask, url_for, redirect, render_template

import requests
api_url = "http://localhost:8080"


# GET
def getAllArticles(url):
    """Get a list with all articles in json format."""
    uri = f'{url}/articles'
    response = requests.get(uri)
    articles = response.json()
    return articles

# DELETE
def deleteArticleById(url, id):
    """Delete an article by Id."""
    uri = f'{url}/article/delete?id={id}'
    response = requests.post(uri)
    return response


# CREATE
def createArticle(url, article):
    """Add new article."""
    uri = f'{url}/article/create'
    response = requests.post(uri, article)
    return response


app = Flask(__name__)

@app.route('/')
def hello():
    return 'Hello, World!'

# endpoint for viewing all articles
@app.route('/articles', methods = ['GET'])
def index():
    articles = getAllArticles(api_url)
    return render_template('index.html', articles=articles)

# endpoint for deleting a record
@app.route('/article/delete?id=<int:article_id>', methods = ['POST'])
def delete_article(article_id):
    print('Deleting article')
    deleteArticleById(api_url, article_id)
    # articles = getAllArticles(api_url)
    # return render_template('index.html', articles=articles)
    return redirect(url_for('index'))

app.run(debug=True)