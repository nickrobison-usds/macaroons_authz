//
// Created by Nicholas Robison on 2019-01-28.
//

#include <QtGui>
#include <boost/interprocess/managed_shared_memory.hpp>
#include <boost/interprocess/sync/interprocess_semaphore.hpp>
#include "LoginComponent.hpp"
#include "SharedMemory.hpp"

static const char *const memory_name = "Macaroons";
namespace Shared {
    namespace bip = boost::interprocess;
    using Segment = bip::managed_shared_memory;
    using Manager = Segment::segment_manager;
    template<typename T>
    using Alloc = bip::allocator<T, Manager>;
}

int main(int argc, char **argv) {



    //Remove shared memory on construction and destruction
    struct shm_remove {
        shm_remove() {
            qDebug() << "Removing";
            Shared::bip::shared_memory_object::remove(memory_name);
        }

        ~shm_remove() {
            qDebug() << "Removing (destructor)";
            Shared::bip::shared_memory_object::remove(memory_name);
        }
    } remover;
    QGuiApplication app(argc, argv);
    LoginComponent comp;

    // Create the shared memory
    Shared::bip::shared_memory_object::remove(memory_name);
    qDebug() << "creating";
    Shared::bip::managed_shared_memory shm(Shared::bip::create_only, memory_name, 65536);
    qDebug() << "created";

    using A = Shared::Alloc<char>;
    A alloc(shm.get_segment_manager());

    auto* data = shm.find_or_construct<SharedMemory<A>>("data")(1024, shm.get_segment_manager());
    qDebug() << "Waiting for Mutex";
    // Wait for the mutex to be ready
    data->mutex.wait();
    qDebug() << "Beginning login";
    comp.login();

    QObject::connect(&comp, &LoginComponent::token, [data](const QString &token_resp) {
        qDebug() << "has token " << token_resp;

        auto token = token_resp.toStdString();
        for (int i=0; i < token.length(); i++) {
            data->mData.at(i) = token[i];
        }
        data->mData.resize(token.length());
        data->done.post();
    });

    QGuiApplication::exec();
}
