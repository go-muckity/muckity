const MongoClient = require('mongodb').MongoClient;
const assert = require('assert');

const dbUser = process.env.MUCKITY_DB_USERNAME || "muckity";
const dbPwd = process.env.MUCKITY_DB_PWD || "muckity";
const dbHost = process.env.MUCKITY_DB_HOST || "localhost";
const dbPort = process.env.MUCKITY_DB_PORT || "27017";
const dbName = process.env.MUCKITY_DB_NAME || "muckity";

// Connection URL
const url = `mongodb://${dbUser}:${dbPwd}@${dbHost}:${dbPort}/${dbName}`;

// Create a new MongoClient
const client = new MongoClient(url);

// Use connect method to connect to the Server
client.connect(function(err) {
    assert.equal(null, err);
    const db = client.db(dbName);

    checkWorldDocument(db, function (){
        client.close();
    });

});

const checkWorldDocument = function(db, callback) {
    const collection = db.collection('worlds');
    collection.find({'_id': 'world:descriptive-aliased-world'}).toArray(function(err, docs) {
        assert.equal(err, null);
        const wdoc = docs[0];
        assert.equal(wdoc["name"], "Descriptive, aliased, world");
        callback(docs)
    });
};