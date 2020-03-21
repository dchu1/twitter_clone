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

  constructor(private router: Router, private apiService: RestApiService) { }

  ngOnInit() {
  }

  signup() {
    let body = {
      FirstName: this.name,
      Lastname: this.lastname,
      Email: this.email,
      Password: this.password
    }
    this.apiService.postData("signup", body).subscribe((response: any) => {
      console.log("[Response]:: ", response);
      this.router.navigate(['./login']);
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
