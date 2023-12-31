import { Module, forwardRef } from '@nestjs/common';
import { UserModule } from 'src/user/user.module';
import { PassportModule } from '@nestjs/passport';
import { NaverStrategy } from './strategy/naver.strategy';

@Module({
  imports: [
    forwardRef(() => UserModule),
    PassportModule.register({
      defaultStrategy: 'jwt',
      session: false,
    }),
  ],
  providers: [NaverStrategy],
  exports: [PassportModule],
})
export class AuthModule {}
