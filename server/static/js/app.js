/**
 * caddy-docker Web UI
 * Copyright 2017 - Nathan Osman
 */

(function() {

    /**
     * Create a <td> element
     */
    function td(t) {
        var isString = $.type(t) === 'string',
            isArray = $.isArray(t),
            $e;
        if (isString || isArray) {
            if (isArray) {
                t = t.join(', ');
            }
            $e = $('<span>').text(t);
        } else {
            $e = t;
        }
        return $('<td>').append($e);
    }

    /**
     * Restart the specified container
     */
    function restartContainer(id) {
        $.post('/api', {action: 'restartContainer', id: id}, function(d) {
            // TODO: error handling
        });
    }

    /**
     * Reloads the list of currently running containers
     */
    function reloadContainers() {
        $.post('/api', {action: 'getContainers'}, function(d) {

            // Sort the containers by name
            d = d.sort(function(a, b) {
                a = a.Name.toUpperCase();
                b = b.Name.toUpperCase();
                if (a < b) {
                    return -1;
                } if (a > b) {
                    return 1;
                } else {
                    return 0;
                }
            });

            // Add the containers to the list
            var $tbody = $('#containers tbody').empty();
            $.each(d, function() {
                var $btns = $('<div>')
                        .addClass('right')
                        .append($('<button>').click(function() {
                            restartContainer(d.ID);
                        }));
                $('<tr>')
                    .append(td(this.Name))
                    .append(td(this.Domains))
                    .append(td(this.Addr))
                    .append(td($btns))
                    .appendTo($tbody);
            });
        });
    }

    reloadContainers();

})();
