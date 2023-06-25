package main

import (
	userservice "api-gateway/userService/kitex_gen/UserService"
	"context"
	"fmt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	userData map[string]*userservice.InsertUser
}

// QueryUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) QueryUser(ctx context.Context, req *userservice.QueryUser) (resp *userservice.QueryUserResponse, err error) {
	// Check if the user ID exists
	fmt.Println("Reached QueryUser")
	user, ok := s.userData[req.ID]
	if !ok {
		return nil, fmt.Errorf("user with ID '%s' not found", req.ID)
	}

	// Print userData map

	fmt.Println("userData:")
	fmt.Println(s.userData)
	fmt.Println("\n Formatted userData:")
	for id, u := range s.userData {
		fmt.Printf("ID: %s, Name: %s, Email: %s, Age: %d\n", id, u.Name, u.Email, u.Age)
	}

	fmt.Println("Reached QueryUser END")

	return &userservice.QueryUserResponse{
		Exist: true,
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

// InsertUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) InsertUser(ctx context.Context, req *userservice.InsertUser) (resp *userservice.InsertUserResponse, err error) {
	// Check if the user ID already exists
	fmt.Println("Reached InsertUser")

	if _, ok := s.userData[req.ID]; ok {
		return &userservice.InsertUserResponse{
			Ok:  false,
			Msg: fmt.Sprintf("User with ID '%s' already exists", req.ID),
		}, nil
	}

	// Insert the user into the map
	if s.userData == nil {
		s.userData = make(map[string]*userservice.InsertUser)
	}

	s.userData[req.ID] = req

	fmt.Println("userData:")
	fmt.Println(s.userData)
	fmt.Println("\n Formatted userData:")
	for id, u := range s.userData {
		fmt.Printf("ID: %s, Name: %s, Email: %s, Age: %d\n", id, u.Name, u.Email, u.Age)
	}

	fmt.Println("Reached InsertUser END")

	return &userservice.InsertUserResponse{
		Ok:  true,
		Msg: fmt.Sprintf("User %s inserted successfully", req.ID),
	}, nil
}
