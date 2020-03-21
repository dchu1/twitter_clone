import { Routes } from '@angular/router';

import { AdminLayoutComponent } from './layouts/admin/admin-layout.component';
import { AuthLayoutComponent } from './layouts/auth/auth-layout.component';

export const AppRoutes: Routes = [
    { path: '', redirectTo: 'login', pathMatch: 'full' },
    { path: '', loadChildren: './login/login.module#LoginModule' },
    { path: 'signup', redirectTo: 'signup', pathMatch: 'full' },
    { path: '', loadChildren: './signup/signup.module#SignupModule' },
    {
        path: '', component: AdminLayoutComponent,
        children: [
            { path: '', loadChildren: './dashboard/dashboard.module#DashboardModule' }
        ]
    },

];
