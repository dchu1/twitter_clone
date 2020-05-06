import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FormGroup, FormControl } from '@angular/forms';
import { RestApiService } from "@sharedComponents/services/rest-api.service";
import { HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  email: string;
  password: string;
  isValid: boolean;
  constructor(private router: Router, private apiService: RestApiService) { }

  ngOnInit() { this.isValid = true }

  login() {

    let body = {
      "Email": this.email,
      "Password": this.password
    }
    this.apiService.postData("login", body).subscribe((response: any) => {
      // console.log("[Response]:: ", response);
      if (response.Status == 200) {
        this.isValid = true
        this.router.navigate(['./home'])
      }
      else if (response.Status == 500){
        console.log("Database server not responding!!")
      }
      else{
        this.isValid = false
        console.log("Login unsuccessful")
      }
    },
      error => {
        console.log("[Error]:: ", error);
      });
  }
  redirectToSignup() {
    this.router.navigate(['./signup'])
  }

}


