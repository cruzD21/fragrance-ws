package crawler

import (
	"bytes"
	"context"
	"fmt"
	"fragrance-ws/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

var bucketName = "fragrance-pictures"

func awsInit() (*s3.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("godotenv.Load -> %v", err)
	}

	var accountId = os.Getenv("ACCOUNT_ID")
	var accessKeyId = os.Getenv("ACCESS_KEY_ID")
	var accessKeySecret = os.Getenv("SECRET_ACCESS_KEY")

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Printf("config.LoadDefaultConfig -> %v", err)
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func putObjectIntoBucket(s3Client *s3.Client, imgID string) error {

	//THIS GETS IMAGE AND PUTS IT
	data, err := getImageByte(imgID)
	if err != nil {
		log.Fatalf("getImageByte -> %v", err)
	}

	fileName := fmt.Sprintf("image-%s.png", imgID)
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		log.Fatalf("Failed to upload object: %v", err)
	}

	log.Println("Successfully uploaded the image to S3")
	return nil
}

func getObject(s3bucket *s3.Client) (*s3.GetObjectOutput, error) {
	object, err := s3bucket.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    aws.String("perfumefavoritodevalentina.png"),
	})
	if err != nil {
		log.Printf("s3bukcet.GetObject -> %v", err)
		return nil, err
	}

	data, err := io.ReadAll(object.Body)
	if err != nil {
		log.Printf("io.ReadAll -> %v", err)
		return nil, err
	}

	err = os.WriteFile("image-return.png", data, 0644)
	if err != nil {
		log.Printf("os.WriteFile -> %v", err)
		return nil, err
	}
	return object, nil
}

//func main() {
//	s3Client := awsClient()
//
//	//listObjectsOutput, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
//	//	Bucket: &bucketName,
//	//})
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//for _, object := range listObjectsOutput.Contents {
//	//	obj, _ := json.MarshalIndent(object, "", "\t")
//	//	fmt.Println(string(obj))
//	//}
//
//	//  {
//	//  	"ChecksumAlgorithm": null,
//	//  	"ETag": "\"eb2b891dc67b81755d2b726d9110af16\"",
//	//  	"Key": "ferriswasm.png",
//	//  	"LastModified": "2022-05-18T17:20:21.67Z",
//	//  	"Owner": null,
//	//  	"Size": 87671,
//	//  	"StorageClass": "STANDARD"
//	//  }
//
//	//listBucketsOutput, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//for _, object := range listBucketsOutput.Buckets {
//	//	obj, _ := json.MarshalIndent(object, "", "\t")
//	//	fmt.Println(string(obj))
//	//}
//
//	// {
//	// 		"CreationDate": "2022-05-18T17:19:59.645Z",
//	// 		"Name": "sdk-example"
//	// }
//}

//	func getR2Image(r2 *s3.Client) []byte {
//		r2.GetObject()
//	}

func getImageByte(imageID string) ([]byte, error) {
	url := fmt.Sprintf("https://fimgs.net/mdimg/perfume/375x500.%s.jpg", imageID)
	res, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get -> %v", err)
		return nil, err
	}
	defer res.Body.Close() // This should be immediately after checking the error from http.Get

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll -> %v", err)
		return nil, err
	}

	return data, nil
}

func getProxyClient() (*http.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	KEY := os.Getenv("KEY")
	HOST := os.Getenv("HOST")
	proxyURL := fmt.Sprintf("http://%s@%s", KEY, HOST)
	return createClient(proxyURL), nil
}

func Test() {
	supa := db.DatabaseConn{}
	err := supa.DatabaseInit()

	client, _ := getProxyClient()
	res, err := CreateRequest(client, "https://www.fragrantica.com/perfume/Puzzle-Parfum/Puzzle-Daylight-91072.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	//printing ip logic

	page, err := parseFragrancePage(res)
	if err != nil {
		log.Printf("error inserting into db with error : %e", err)
	}
	fmt.Printf("%+v\n", page)

	err = supa.InsertPage(page)
	if err != nil {
		log.Printf("error inserting into db with error : %e", err)
	}

	log.Println("code inserted  page successfully to db")
}
