import { Component } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Store } from '@ngxs/store';
import { Broadcast } from 'app/model/broadcast.model';
import { NavbarProjectData } from 'app/model/navbar.model';
import { User } from 'app/model/user.model';
import { BroadcastService } from 'app/service/broadcast/broadcast.service';
import { BroadcastStore } from 'app/service/broadcast/broadcast.store';
import { NavbarService } from 'app/service/navbar/navbar.service';
import { PathItem } from 'app/shared/breadcrumb/breadcrumb.component';
import { AutoUnsubscribe } from 'app/shared/decorator/autoUnsubscribe';
import { ToastService } from 'app/shared/toast/ToastService';
import { AuthenticationState } from 'app/store/authentication.state';
import { Subscription } from 'rxjs';
import { finalize } from 'rxjs/operators';

@Component({
    selector: 'app-broadcast-add',
    templateUrl: './broadcast.add.html',
    styleUrls: ['./broadcast.add.scss']
})
@AutoUnsubscribe()
export class BroadcastAddComponent {
    loading = false;
    deleteLoading = false;
    broadcast: Broadcast;
    currentUser: User;
    canAdd = false;
    broadcastLevelsList: any;
    projects: Array<NavbarProjectData> = [];
    navbarSub: Subscription;
    path: Array<PathItem>;

    constructor(
        private _navbarService: NavbarService,
        private _broadcastStore: BroadcastStore,
        private _toast: ToastService, private _translate: TranslateService,
        private _route: ActivatedRoute, private _router: Router,
        private _store: Store,
        private _broadcastService: BroadcastService
    ) {
        this.navbarSub = this._navbarService.getData(true).subscribe((data) => {
            this.loading = false;
            if (Array.isArray(data)) {
                let voidProj = new NavbarProjectData();
                voidProj.type = 'project';
                voidProj.name = ' ';
                this.projects = [voidProj].concat(data.filter((elt) => elt.type === 'project'));
                this.currentUser = this._store.selectSnapshot(AuthenticationState.user);
                this.broadcastLevelsList = this._broadcastService.getBroadcastLevels();
            }
        });

        this._route.params.subscribe(params => {
            this.broadcast = new Broadcast();
        });

        this.path = [<PathItem>{
            translate: 'common_admin'
        }, <PathItem>{
            translate: 'broadcast_list_title',
            routerLink: ['/', 'admin', 'broadcast']
        }, <PathItem>{
            translate: 'common_create'
        }];
    }

    clickSaveButton(): void {
        if (!this.broadcast.title) {
            return;
        }

        this.loading = true;
        this._broadcastStore.create(this.broadcast)
            .pipe(finalize(() => this.loading = false))
            .subscribe(bc => {
                this._toast.success('', this._translate.instant('broadcast_saved'));
                this._router.navigate(['admin', 'broadcast', bc.id]);
            });
    }
}
