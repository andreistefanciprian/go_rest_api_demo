

```
# create python3 virtual env
python3 -m venv .env

# activate python env
source .env/bin/activate

# install pip packages
pip install -r requirements.txt

# deactivate python env
deactivate

export FLASK_APP=app
flask run


docker image build -t frontend:latest .
docker container run --name frontend -p 8090:5000 frontend:latest
```