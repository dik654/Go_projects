"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.UserController = void 0;
const common_1 = require("@nestjs/common");
const naver_auth_guard_1 = require("../auth/guard/naver-auth.guard");
let UserController = class UserController {
    async naverlogin() {
        return;
    }
    async callback(req, res) {
        res.cookie('access_token', req.user.accessToken);
        res.redirect('http://localhost:8080/qrcode');
        res.end();
    }
};
__decorate([
    (0, common_1.UseGuards)(naver_auth_guard_1.NaverAuthGuard),
    (0, common_1.Get)('auth/naver'),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", []),
    __metadata("design:returntype", Promise)
], UserController.prototype, "naverlogin", null);
__decorate([
    (0, common_1.UseGuards)(naver_auth_guard_1.NaverAuthGuard),
    (0, common_1.Get)('auth/nlogincallback'),
    __param(0, (0, common_1.Req)()),
    __param(1, (0, common_1.Res)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object, Object]),
    __metadata("design:returntype", Promise)
], UserController.prototype, "callback", null);
UserController = __decorate([
    (0, common_1.Controller)('user')
], UserController);
exports.UserController = UserController;
//# sourceMappingURL=user.controller.js.map