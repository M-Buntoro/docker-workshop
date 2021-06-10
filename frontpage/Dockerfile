FROM node:4-slim
RUN npm install express
RUN npm install redis
COPY . .
CMD ["node", "frontpage.js"]
EXPOSE 11980
