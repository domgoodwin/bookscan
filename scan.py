import os, requests
from datetime import datetime
import pyttsx3
engine = pyttsx3.init()



auth_token = os.getenv("AUTH_TOKEN")
out_file = "books.csv"

class Book:

    def lookup(self, isbn):
        found, json = lookup_isbn_openlibrary(isbn)
        if found:
            self.import_openlibrary(isbn, json)
    
    def import_openlibrary(self, isbn, json):
        self.title = json['title']
        self.author = lookup_author_openlibrary(json['authors'][0]['key'])
        self.isbn = isbn
        self.pages = json['number_of_pages']
        self.link = "https://openlibrary.org" + json['key']

    def info(self):
        info  = "{0} by {1}".format(self.title, self.author)
        engine.say(info)
        engine.runAndWait()
        print(info)
    
    def store_csv(self):
        # Less efficient opening and closing file each time but means we don't lose book info if program errors
        file = open(out_file, 'a')
        file.write("{},{},{},{},{}".format(self.title, self.author, self.isbn, self.pages, self.link))
        file.close()

def lookup_isbn_google(isbn):
    url = "https://www.googleapis.com/books/v1/"
    search_path = "volumes?q="
    try: 
        rsp = requests.get("{0}{1}isbn:{2}".format(url, search_path,isbn))
        if rsp.status_code != 200:
            print("non 200: {0}, {1}".format(rsp.status_code, rsp.content))
        json_rsp = rsp.json()
        if "items" in json_rsp:
            book = json_rsp['items'][0]['volumeInfo']
            print("{0} by {1}".format(book['title'], book['authors'][0]))
            return True, json_rsp # TODO: Return book struct
        return False
    except Exception as err:
        print("error occured: {0}".format(err))

def lookup_isbn_openlibrary(isbn):
    url = "https://openlibrary.org/isbn/{0}.json".format(isbn)
    try: 
        rsp = requests.get(url)
        if rsp.status_code != 200:
            print("non 200: {0}, {1}".format(rsp.status_code, rsp.content))
        return True, rsp.json()
    except Exception as err:
        print("error occured: {0}".format(err))

def lookup_author_openlibrary(author_id):
    url = "https://openlibrary.org{0}.json".format(author_id)
    try: 
        rsp = requests.get(url)
        if rsp.status_code != 200:
            print("non 200: {0}, {1}".format(rsp.status_code, rsp.content))
        return rsp.json()['name']
    except Exception as err:
        print("error occured: {0}".format(err))


while True:
    isbn = input("scan book... ")
    book = Book()
    book.lookup(isbn)
    book.info()
    book.store_csv()

