import { Injectable } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';

@Injectable()
export class AuthService {
    constructor(
        private jwtService: JwtService,
    ){}
    onceToken(user_profile: any) {
        const payload = {
          user_email: user_profile.user_email,
          user_nick: user_profile.user_nick,
          user_provider: user_profile.user_provider,
          user_token: 'onceToken',
        };
    
        return this.jwtService.sign(payload, {
          secret: process.env.JWT_SECRET,
          expiresIn: '10m',
        });
      }
}
