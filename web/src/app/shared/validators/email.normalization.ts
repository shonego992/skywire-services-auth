 export class EmailNormalization {
       public emailNormalization(emailValue: string) {
        if (emailValue.includes('@gmail.com')) {
            var parts = emailValue.split('@');
            parts[0] = parts[0].split('.').join('');
            emailValue = parts[0] + '@' +  parts[1];
          }
          return emailValue;
    }
}