import { Component, OnInit } from "@angular/core";
import { RestApiService } from "@sharedComponents/services/rest-api.service";
import { Router } from "@angular/router";

@Component({
  selector: "app-dashboard",
  templateUrl: "./dashboard.component.html",
  styleUrls: ['./dashboard.component.scss']
})

export class DashboardComponent implements OnInit {
  newTweet: string;
  newsData: Array<any>;
  usersData: Array<any>;

  constructor(private router: Router, private apiService: RestApiService) { }

  ngOnInit() {
    this.newsData = [{ news: "News 1" }, { news: "News 2" }, { news: "News 3" }, { news: "News 4" }];
    this.usersData = [{ username: "User 1", status: "Follow" }, { username: "User 2", status: "Follow" }, { username: "User 3", status: "Follow" }, { username: "User 4", status: "Follow" }];
  }

  logout() {
    this.apiService.getData("logout").subscribe((response: any) => {
      console.log("[Response]:: ", response);
      this.router.navigate(['./login']);
    },
      error => {
        console.log("[Error]:: ", error);
        this.router.navigate(['./login']);
      });
    // this.router.navigate(['./home']);

  }

  postUpdate() {
    this.newTweet = "";
  }

  userStatus(i) {
    this.usersData.forEach((item, index) => {
      if (index == i) {
        if (item.status == "Follow")
          item.status = "Unfollow";
        else
          item.status = "Follow"
      }
    });
  }

  /**Use this inside a function
   * login(){
   * api call given below
   * }
    * Example of Get api
    *
    this.apiService.getData("newsFeed").subscribe((response: any) => {
      console.log("[Response]:: ", response);
        perform any line of code here
        ex: this.router.navigate(['./login']);
    },
      error => {
        console.log("[Error]:: ", error)
      });
  }
  */

  /**
    * Function of Post api
    *
    let body = {
      "email": "test user",
      "password": "sajhdg"
    }
    this.apiService.postData("login", body).subscribe((response: any) => {
      console.log("[Response]:: ", response);
        code here
    },
      error => {
        console.log("[Error]:: ", error);
      });
  }
  */
}
