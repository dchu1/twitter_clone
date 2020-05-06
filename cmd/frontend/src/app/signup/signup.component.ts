import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { RestApiService } from "@sharedComponents/services/rest-api.service";

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})
export class SignupComponent implements OnInit {
  name: string;
  lastname: string;
  email: string;
  password: string;
  isValidUsername: boolean;

  constructor(private router: Router, private apiService: RestApiService) { }

  ngOnInit() {
    this.isValidUsername = true
  }

  signup() {
    let body = {
      FirstName: this.name,
      Lastname: this.lastname,
      Email: this.email,
      Password: this.password
    }
    this.apiService.postData("signup", body).subscribe((response: any) => {
      // console.log("[Response]:: ", response);
       if (response.Status == 201) {
        this.isValidUsername = true
        this.router.navigate(['./login']);
      }
      else if (response.Status == 500){
        console.log("Database server not responding!!")
      }
      else {
        this.isValidUsername = false
        console.log("Sign Up unsuccessful")
      }
      // code here
    },
      error => {
        console.log("[Error]:: ", error);
        this.router.navigate(['./signup']);
      });
    // this.router.navigate(['./home']);
  }

  redirectToLogin() {
    this.router.navigate(['./login']);
  }
}
