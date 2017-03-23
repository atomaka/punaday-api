const axios = require('axios')

const PUN_BASE = 'http://www.punoftheday.com'

module.exports.today = ((event, context, callback) => {
  respondWithPunFrom(PUN_BASE, callback)
})

module.exports.random = ((event, context, callback) => {
  respondWithPunFrom(`${PUN_BASE}/cgi-bin/randompun.pl`, callback)
})

module.exports.show = ((event, context, callback) => {
  const id = event.pathParameters.id

  respondWithPunFrom(`${PUN_BASE}/pun/${id}`, callback)
})

function respondWithPunFrom(url, callback) {
  return axios.get(url)
    .then(response => {
      const pun = parsePun(response.data)

      callback(null, punResponse(pun))
    })
    .catch(error => {
      callback(null, errorResponse())
    })
}

function parsePun(html) {
  const punMatches = html.match(/<p>(.*)<\/p>/)
  const text = punMatches[1].replace('&#8220;', '').replace('&#8221;', '')
  const idMatches = html.match(/name="PunID" value="(\d+)"/)
  const id = idMatches[1]
  const urlMatches = html.match(/class="fb-share-button" data-href="(.*)" d/)
  const url = urlMatches[1]

  return { id, text, url }
}

function punResponse (pun) {
  delete pun['id']

  return {
    statusCode: 200,
    body: JSON.stringify(pun)
  }
}

function errorResponse () {
  return { statusCode: 500 }
}
