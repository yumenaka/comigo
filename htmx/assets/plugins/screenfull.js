import screenfull from 'screenfull'

if(document.getElementById('FullScreenIcon')){
    document.getElementById('FullScreenIcon').addEventListener('click', () => {
        if (screenfull.isEnabled) {
            screenfull.toggle()
        } else {
            // Ignore or do something else
            i18next.t('not_support_fullscreen')
        }
    })
} 