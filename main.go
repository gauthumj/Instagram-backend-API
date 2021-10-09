package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Input struct {
    Name string `json:"name"`
    Email string `json:"email"`
    Pass string `json:"pass"`
}

type Post struct {
    Caption string `json:"caption"`
    ImgaeUrl string `json:"imageUrl"`
    Userid primitive.ObjectID `json:"userid"`
}


func Encrypt(key, data []byte) ([]byte, error) {
    blockCipher, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(blockCipher)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = rand.Read(nonce); err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nonce, nonce, data, nil)

    return ciphertext, nil
}

func Decrypt(key, data []byte) ([]byte, error) {
    blockCipher, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(blockCipher)
    if err != nil {
        return nil, err
    }

    nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}

func GenerateKey() ([]byte, error) {
    key := make([]byte, 32)

    _, err := rand.Read(key)
    if err != nil {
        return nil, err
    }

    return key, nil
}


func writeToDatabase(user Input) {
    client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("KEY")))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    db := client.Database("Insta")
    collection := db.Collection("users")
    var (
        
        data     = []byte(user.Pass)
    )
    key, err := GenerateKey()
    if err != nil {
        log.Fatal(err)
    }
    encPass, err := Encrypt(key, data)  
    if err != nil {
        log.Fatal(err)
    }

    
    collection.InsertOne(ctx,bson.D{
        {Key: "Name", Value: user.Name},
        {Key: "Email", Value: user.Email},
        {Key: "Pass", Value: hex.EncodeToString(encPass)},
    })

}

func createUser(w http.ResponseWriter, r *http.Request) {    
    var input Input
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        w.WriteHeader(400)
        fmt.Fprintf(w, "Decode error! please check your JSON formating.")
        return
    }

    writeToDatabase(input)
    fmt.Fprintln(w, "User created successfully!")
}

func getUser(w http.ResponseWriter, r *http.Request) {
    id1 := r.URL.Query().Get("id")
    id, _ := primitive.ObjectIDFromHex(id1)
    
    client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("KEY")))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    db := client.Database("Insta")
    collection := db.Collection("users")
    filterCursor, err := collection.Find(ctx, bson.M{"_id": id}) 
    if err != nil {
        log.Fatal(err)
    }
    var episodesFiltered []bson.M
    if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
        log.Fatal(err)
    }
    fmt.Fprintln(w, episodesFiltered)
}

func createPost(w http.ResponseWriter, r *http.Request) {
    
    var input Post
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        w.WriteHeader(400)
        fmt.Fprintf(w, "Decode error! please check your JSON formating.")
        return
    }

    client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("KEY")))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    db := client.Database("Insta")
    collection := db.Collection("posts")
    collection.InsertOne(ctx,bson.D{
        {Key: "Caption", Value: input.Caption},
        {Key: "URL", Value: input.ImgaeUrl},
        {Key: "Timestamp", Value: time.Now().String()},
        {Key: "Userid", Value: input.Userid},
    })
    fmt.Fprintln(w, "Post created")
}

func getPosts(w http.ResponseWriter, r *http.Request) {
    id1 := r.URL.Query().Get("id")
    id, _ := primitive.ObjectIDFromHex(id1)
    print(os.Getenv("KEY"))
    client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("KEY")))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    db := client.Database("Insta")
    collection := db.Collection("posts")
    filterCursor, err := collection.Find(ctx, bson.M{"_id": id}) 
    if err != nil {
        log.Fatal(err)
    }
    var episodesFiltered []bson.M
    if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
        log.Fatal(err)
    }
    fmt.Fprintln(w, episodesFiltered[0])
}

func userPosts(w http.ResponseWriter, r *http.Request){
    id1 := r.URL.Query().Get("id")
    id, _ := primitive.ObjectIDFromHex(id1)
    lim := r.URL.Query().Get("lim")
    number,err := strconv.ParseUint(lim, 10, 32)
    finalLim := int(number) //Convert uint64 To int
    
    client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("KEY")))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    db := client.Database("Insta")
    collection := db.Collection("posts")
    filterCursor, err := collection.Find(ctx, bson.M{"Userid": id}) 
    if err != nil {
        log.Fatal(err)
    }
    var episodesFiltered []bson.M
    if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
        log.Fatal(err)
    }
    fmt.Fprintln(w, episodesFiltered[:(finalLim)])
}


func main() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
    http.HandleFunc("/users", createUser)
    http.HandleFunc("/users/", getUser)
    http.HandleFunc("/posts", createPost)
    http.HandleFunc("/posts/", getPosts)
    http.HandleFunc("/posts/users/", userPosts)
    log.Fatal(http.ListenAndServe(":8080", nil))
}