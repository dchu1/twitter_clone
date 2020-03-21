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

  constructor(private router: Router, private apiService: RestApiService) { }

  ngOnInit() { }

  login() {

    let body = {
      "Email": this.email,
      "Password": this.password
    }
    this.apiService.postData("login", body).subscribe((response: any) => {
      this.router.navigate(['./home'])
    },
      error => {
        console.log("[Error]:: ", error);
      });
  }
  redirectToSignup() {
    this.router.navigate(['./signup'])
  }

}


