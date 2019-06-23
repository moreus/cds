import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthenticationService } from 'app/service/services.module';
import { finalize } from 'rxjs/operators';

@Component({
    selector: 'app-auth-verify',
    templateUrl: './verify.html',
    styleUrls: ['./verify.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class VerifyComponent implements OnInit {
    userVerified: any;
    userVerifiedDisplay: any;
    showErrorMessage: boolean;
    loading: boolean;

    constructor(
        private _router: Router,
        private _route: ActivatedRoute,
        private _authenticationService: AuthenticationService,
        private _cd: ChangeDetectorRef
    ) {
        this.loading = true;
    }

    ngOnInit(): void {
        this._route.queryParams.subscribe(queryParams => {
            let token = queryParams['token'];
            if (!token) {
                this.showErrorMessage = true;
                this.loading = false;
                this._cd.detectChanges();
                return;
            }

            this._authenticationService.localVerify(token)
                .pipe(finalize(() => {
                    this.loading = false;
                    this._cd.detectChanges();
                }))
                .subscribe(res => {
                    // TODO show successfull verify with redirect button
                    // this._router.navigate(['home']);
                }, () => {
                    this.showErrorMessage = true;
                });
        });
    }

    navigateToSignin() {
        this._router.navigate(['/auth/signin']);
    }
}
