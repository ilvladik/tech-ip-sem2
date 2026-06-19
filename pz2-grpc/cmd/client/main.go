package main

import (
	"context"
	"log"
	"time"

	"example.com/pz2-grpc/gen/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := studentpb.NewStudentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("=== Ping ===")
	pingResp, err := client.Ping(ctx, &studentpb.PingRequest{Message: "hello grpc"})
	if err != nil {
		log.Fatal("Ping error:", err)
	}
	log.Println("Ping response:", pingResp.GetMessage())

	log.Println("=== GetStudentByID id=1 ===")
	studentResp, err := client.GetStudentByID(ctx, &studentpb.GetStudentRequest{Id: 1})
	if err != nil {
		log.Fatal("GetStudentByID error:", err)
	}
	printStudent(studentResp.GetStudent())

	log.Println("=== ListStudents ===")
	listResp, err := client.ListStudents(ctx, &studentpb.Empty{})
	if err != nil {
		log.Fatal("ListStudents error:", err)
	}
	log.Printf("Total students: %d", len(listResp.GetStudents()))
	for _, st := range listResp.GetStudents() {
		printStudent(st)
	}

	log.Println("=== CreateStudent ===")
	created, err := client.CreateStudent(ctx, &studentpb.CreateStudentRequest{
		FullName:       "Козлов Дмитрий Андреевич",
		Group:          "ИТТ-04-25",
		Email:          "kozlov@example.com",
		Specialization: "Машинное обучение",
	})
	if err != nil {
		log.Fatal("CreateStudent error:", err)
	}
	printStudent(created.GetStudent())

	log.Println("=== GetStudentByID new id ===")
	studentResp, err = client.GetStudentByID(ctx, &studentpb.GetStudentRequest{Id: created.GetStudent().GetId()})
	if err != nil {
		log.Fatal("GetStudentByID error:", err)
	}
	printStudent(studentResp.GetStudent())

	log.Println("=== GetStudentByID id=999 (NotFound) ===")
	_, err = client.GetStudentByID(ctx, &studentpb.GetStudentRequest{Id: 999})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Println("Expected NotFound:", err)
		} else {
			log.Fatal("unexpected error:", err)
		}
	} else {
		log.Fatal("expected NotFound error, got success")
	}
}

func printStudent(st *studentpb.Student) {
	log.Printf("  id=%d, name=%s, group=%s, email=%s, spec=%s",
		st.GetId(), st.GetFullName(), st.GetGroup(), st.GetEmail(), st.GetSpecialization())
}
