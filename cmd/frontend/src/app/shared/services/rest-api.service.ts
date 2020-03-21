import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { environment } from '@environments/environment';
import { Observable } from "rxjs";

@Injectable()
export class RestApiService {
  apiURL: string = environment.backendUrl;

  constructor(private httpClient: HttpClient) { }

  /**
   * Function for http call to get api
   * @param: url
   */
  getData = (endpoint: string) => {
    return this.httpClient.get(this.apiURL + endpoint,{
      withCredentials: true,
    });
  };
  /**
   * Function for http call to post api
   * @param endpoint: url
   * @param body: data to send
   */
  postData = (endpoint: string, body: object, ):Observable<any> => {
    return this.httpClient.post(this.apiURL + endpoint, body,{
      withCredentials: true,
    });
  }
}
