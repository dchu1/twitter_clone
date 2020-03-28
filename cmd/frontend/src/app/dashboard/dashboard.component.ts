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
    // this.newsData = [{ news: "News 1" }, { news: "News 2" }, { news: "News 3" }, { news: "News 4" }];
    // this.usersData = [{ username: "User 1", status: "Follow" }, { username: "User 2", status: "Follow" }, { username: "User 3", status: "Follow" }, { username: "User 4", status: "Follow" }];

    this.getNewsFeed();
    this.getUserList();
  }

  getNewsFeed() {
    this.apiService.getData("feed").subscribe((response: any) => {
      console.log("[Response]:: ", response);
      this.newsData = []
      if (response["posts"])
        this.newsData = response.posts
    },
      error => {
        console.log("[Error]:: ", error);
      });

  }

  getUserList() {
    this.usersData = []
    this.apiService.getData("user/following").subscribe((response: any) => {
      console.log("[Response]  followers::", response);
      let userArray = response["Users"]
      if (userArray.length > 0) {
        userArray.forEach(element => {
          this.usersData.push({ firstName: element.firstname, lastName: element.lastname, userId: element.userId, status: "Unfollow" })
        });
      }
    },
      error => {
        console.log("[Error]:: ", error);
      });
    this.apiService.getData("user/notfollowing").subscribe((response: any) => {
      console.log("[Response] not followers:: ", response);
      let userArray = response["Users"]
      if (userArray.length > 0) {
        userArray.forEach(element => {
          this.usersData.push({ firstName: element.firstname, lastName: element.lastname, userId: element.userId, status: "Follow" })
        });
      }
    },
      error => {
        console.log("[Error]:: ", error);
      });

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

  postTweet() {
    let body = {
      "Message": this.newTweet,
    }
    this.apiService.postData("post", body).subscribe((response: any) => {
      console.log("[Response]:: ", response);
      this.router.navigate(['./home']);
      this.newTweet = ""
    },
      error => {
        console.log("[Error]:: ", error);
        this.router.navigate(['./login']);
      });

  }

  userStatus(i) {
    this.usersData.forEach((item, index) => {
      if (index == i) {
        let body = {
          "UserId": item.userId,
        }
        if (item.status == "Follow") {
          this.apiService.postData("follow/create", body).subscribe((response: any) => {
            console.log("[Response]:: ", response);
        
          },
            error => {
              console.log("[Error]:: ", error);
            });
          item.status = "Unfollow";
        }
        else {
          let body = {
            "UserId": item.userId,
          }
          this.apiService.postData("follow/destroy", body).subscribe((response: any) => {
            console.log("[Response]:: ", response);
          },
            error => {
              console.log("[Error]:: ", error);
            });
          item.status = "Follow"
        }

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
