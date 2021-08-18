conn = new Mongo();
db = conn.getDB("resort");

db = connect("localhost:27017/resort");

db.users.drop();

db.createCollection("users", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["fisrtname", "lastname", "email", "password", "phone_number"],
            properties: {
                firstname: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                lastname: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                email: {
                    unique: true,
                    bsonType: "string",
                    // Regexp to validate emails with more strict rules as added in tests/users.js which also conforms mostly with RFC2822 guide lines
                    match: [/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/, 'Please enter a valid email'],
                },
                password: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                phone_number: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                hotel: {
                    bsonType: "array",
                    description: "reserved hotel rooms"
                },
                restaurant: {
                    bsonType: "array",
                    description: "ordered foods"
                },
                profile: {
                    bsonType: "array",
                    description: "more personal information"
                },

            }
        }
    }
});

db.users.insertOne({ firstname: "arian", lastname: "pourarian", email: "ahmadgholi@gmail.com", password: "123", phone_number: "00989054778974" });