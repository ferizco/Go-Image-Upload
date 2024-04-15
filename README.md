# Go Image Upload 
Web Application to store image built using Go Fiber and Postgresql with security approach 


## Feature 
- Signup, signin, & signout
- Add image & delete image

## Tech 
- Golang Fiber v2.52.4
- PostgreSQL v1.5.7
- Gorm v1.25.9
- Bcrypt
- JWT
- Aes
- Sweet Alert

## Security 
- JWT cookie encrypted
- CSRF token
- Password hashing
- Strong password filtering
- Rate limiter
- Session Expiring
- Image Validation (Format & oversize limiter)
- Logging

### Flow Chart 
![Code drawio](https://github.com/ferizco/Go-Image-Uploadv2/assets/71595406/c4f1d3ae-4eb9-4125-9487-d30000cf6719)

### How to Start 
1. Install GO and Postgresql
2. in database.go file on dsn variable set up the db configuration (user=USER_DB password=PASSWORD_USER_DB dbname=DB_NAME port=DB_PORT)
3. run the go server in terminal
4. open localhost on port 5000

### Functional Test 
| API | Description | Status |
| ----- | ---------- | -----|
| / | index page | Passed |
| /api/user/login | for login | Passed |
| /api/user/dashboard | show dashboard | Passed |
| /api/user/logout | for logout | Passed |
| /api/user/signup | for signup | Passed |
| /api/image/upload | for upload image | Passed |
| /api/image/delete-image/:id | for delete image | Passed |

for evidence of testing please follow the link below

[![Youtube](https://img.youtube.com/vi/ox77gO_zs0Q/0.jpg)](https://www.youtube.com/watch?v=ox77gO_zs0Q)


### User Interface
<img src="https://github.com/ferizco/Go-Image-Uploadv2/assets/71595406/d8e8f933-06dc-4861-99dd-6a7c6abca3f7" alt="Alt Text" width="500">
<img src="https://github.com/ferizco/Go-Image-Uploadv2/assets/71595406/6f800809-821c-4375-a0ee-40ce00a4c532" alt="Alt Text" width="500">
<img src="https://github.com/ferizco/Go-Image-Uploadv2/assets/71595406/76abb8de-2c16-4a09-ba17-eb5cc03d3e2e" alt="Alt Text" width="500">




