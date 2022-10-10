from flask import Flask, url_for, redirect, render_template
import requests
import json
from flask_bootstrap import Bootstrap5
from flask_wtf import FlaskForm
from wtforms import StringField, SubmitField, TextAreaField
from wtforms.validators import DataRequired
import os

REST_API_HOST=os.environ.get('REST_API_HOST')
REST_API_PORT=os.environ.get('REST_API_PORT')
api_url = f'http://{REST_API_HOST}:{REST_API_PORT}'


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

# UPDATE
def updateArticleById(url, id, article):
    """Update an article by Id."""
    uri = f'{url}/article/update?id={id}'
    response = requests.post(uri, article)
    return response

# CREATE
def createArticle(url, article):
    """Add new article."""
    uri = f'{url}/article/create'
    response = requests.post(uri, article)
    return response


app = Flask(__name__)

# Flask-WTF requires an encryption key - the string can be anything
app.config['SECRET_KEY'] = 'C2HWGVoMGfNTBsrYQg8EcMrdTimkZfAb'

# Flask-Bootstrap requires this line
bootstrap = Bootstrap5(app)

@app.route('/')
def hello():
    # return "Welcome to the Library"
    return redirect(url_for('articles'))

# endpoint for viewing all records
@app.route('/articles', methods = ['GET'])
def articles():
    articles = getAllArticles(api_url)
    return render_template('articles.html', articles=articles)

# endpoint for deleting a record
@app.route('/article/delete?id=<int:article_id>', methods = ['POST'])
def delete_article(article_id):
    deleteArticleById(api_url, article_id)
    return redirect(url_for('articles'))

class AddArticle(FlaskForm):
    title = StringField('Book Title', validators=[DataRequired()])
    description = StringField('Book Description', validators=[DataRequired()])
    content = TextAreaField('Book Content', validators=[DataRequired()])
    submit = SubmitField('Submit', validators=[DataRequired()])

# endpoint for adding a record
@app.route('/article/create', methods = ['GET', 'POST'])
def add_article():
    
    form = AddArticle()

    if form.validate_on_submit():
        title = form.title.data
        description = form.description.data
        content = form.content.data

        article = { "Title": title, "desc": description, "content": content }
        article = json.dumps(article)
        createArticle(api_url, article)
        return redirect(url_for('articles'))
    return render_template('add.html',form=form)


class UpdateArticle(FlaskForm):
    id = StringField('Book ID', validators=[DataRequired()])
    title = StringField('Book Title', validators=[DataRequired()])
    description = StringField('Book Description', validators=[DataRequired()])
    content = TextAreaField('Book Content', validators=[DataRequired()])
    submit = SubmitField('Submit', validators=[DataRequired()])

# endpoint for updating a record
@app.route('/article/update', methods = ['GET', 'POST'])
def update_article():
    
    form = UpdateArticle()

    if form.validate_on_submit():
        id = form.id.data
        title = form.title.data
        description = form.description.data
        content = form.content.data

        article = { "Title": title, "desc": description, "content": content }
        article = json.dumps(article)
        updateArticleById(api_url, id, article)
        return redirect(url_for('articles'))
    return render_template('update.html',form=form)

if __name__ == "__main__":
    app.run(debug=True)
    port = int(os.environ.get('PORT', 5000))
    app.run(host='0.0.0.0', debug=True, port=port)