import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { SignupComponent } from './signup.component';
import { SignupRoutes } from './signup.routing';
import { RestApiService } from '@components/shared/services/rest-api.service';

@NgModule({
  declarations: [SignupComponent],
  imports: [
    CommonModule,
    RouterModule.forChild(SignupRoutes),
    FormsModule
  ],
  providers:[RestApiService]
})
export class SignupModule { }
